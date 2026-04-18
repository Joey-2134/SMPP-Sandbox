package dev.joey.smppclient.pdu;

import java.util.Arrays;

import lombok.AllArgsConstructor;
import lombok.Data;
import lombok.NoArgsConstructor;

@Data
@AllArgsConstructor
@NoArgsConstructor
public class UnbindResp {
    private Header header;

    public static UnbindResp fromBytes(byte[] data) {
    if (data.length < Header.LENGTH) {
        throw new IllegalArgumentException("data too short to contain a unbind_resp PDU");
    }
    byte[] headerBytes = Arrays.copyOfRange(data, 0, Header.LENGTH);
    Header header = Header.fromBytes(headerBytes);
    return new UnbindResp(header);
    }
}
