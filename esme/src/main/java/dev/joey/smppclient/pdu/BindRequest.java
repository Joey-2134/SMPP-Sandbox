package dev.joey.smppclient.pdu;

import java.io.ByteArrayOutputStream;
import java.io.IOException;

public class BindRequest {
    private final int commandId;
    private final int sequenceNumber;
    private final String systemId;
    private final String password;
    private final String systemType;

    public BindRequest(int commandId, int sequenceNumber, String systemId, String password, String systemType) {
        this.commandId = commandId;
        this.sequenceNumber = sequenceNumber;
        this.systemId = systemId;
        this.password = password;
        this.systemType = systemType;
    }

    public byte[] toBytes() {
        byte[] systemIdBytes = Utils.toCOctetString(systemId);
        byte[] passwordBytes = Utils.toCOctetString(password);
        byte[] systemTypeBytes = Utils.toCOctetString(systemType);
        byte[] addressRangeBytes = Utils.toCOctetString("");

        int bodyLength = systemIdBytes.length + passwordBytes.length + systemTypeBytes.length + addressRangeBytes.length + 3;
        Header header = new Header(Header.LENGTH + bodyLength, commandId, 0, sequenceNumber);

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
