package dev.joey.smppclient.pdu;

import java.io.ByteArrayOutputStream;
import java.io.IOException;

import lombok.AllArgsConstructor;
import lombok.Data;
import lombok.NoArgsConstructor;

@Data
@AllArgsConstructor
@NoArgsConstructor
public class BindTransmitter {
    private Header header;
    private String systemId;
    private String password;
    private String systemType;
    private int interfaceVersion;
    private int addrTon;
    private int addrNpi;
    private String addressRange;

    public byte[] toBytes() {
        byte[] systemIdBytes = Utils.toCOctetString(systemId);
        byte[] passwordBytes = Utils.toCOctetString(password);
        byte[] systemTypeBytes = Utils.toCOctetString(systemType);
        byte[] addressRangeBytes = Utils.toCOctetString(addressRange);

        int bodyLength = systemIdBytes.length + passwordBytes.length + systemTypeBytes.length + addressRangeBytes.length + 3;
        int commandLength = Header.LENGTH + bodyLength;

        Header header = new Header(commandLength, CommandId.BIND_TRANSMITTER, 0, this.header.getSequenceNumber());

        ByteArrayOutputStream out = new ByteArrayOutputStream();
        try {
            out.write(header.toBytes());
            out.write(systemIdBytes);
            out.write(passwordBytes);
            out.write(systemTypeBytes);
            out.write(Utils.INTERFACE_VERSION);
            out.write(0x00);
            out.write(0x00);
            out.write(addressRangeBytes);
        } catch (IOException e) {
            throw new RuntimeException(e);
        }
        return out.toByteArray();
    }
}
