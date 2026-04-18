package dev.joey.smppclient.pdu;
import java.util.Arrays;

import lombok.AllArgsConstructor;
import lombok.Data;
import lombok.NoArgsConstructor;

@Data
@AllArgsConstructor
@NoArgsConstructor
public class EnquireLinkResp {
    private Header header;

    public static EnquireLinkResp fromBytes(byte[] data) {
        if (data.length < Header.LENGTH) {
            throw new IllegalArgumentException("data too short to contain a enquire_link_resp PDU");
        }
        byte[] headerBytes = Arrays.copyOfRange(data, 0, Header.LENGTH);
        Header header = Header.fromBytes(headerBytes);
        return new EnquireLinkResp(header);
    }

    public byte[] toBytes() {
        return header.toBytes();
    }
}
