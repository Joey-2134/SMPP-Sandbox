package dev.joey.smppclient.pdu;

import java.util.Arrays;

import lombok.AllArgsConstructor;
import lombok.Data;
import lombok.NoArgsConstructor;

@Data
@AllArgsConstructor
@NoArgsConstructor
public class BindTransmitterResp {
    private Header header;
    private String systemId;

    public static BindTransmitterResp fromBytes(byte[] data) {
        if (data.length < Header.LENGTH) {
            throw new IllegalArgumentException("data too short to contain a bind_transmitter_resp PDU");
        }
        byte[] headerBytes = Arrays.copyOfRange(data, 0, Header.LENGTH);
        Header header = Header.fromBytes(headerBytes);
        String systemId = Utils.fromCOctetString(Arrays.copyOfRange(data, Header.LENGTH, data.length));
        return new BindTransmitterResp(header, systemId);
    }
}
