package dev.joey.smppclient.pdu;

import java.util.Arrays;

import lombok.AllArgsConstructor;
import lombok.Data;
import lombok.NoArgsConstructor;

@Data
@AllArgsConstructor
@NoArgsConstructor
public class SubmitSmResp {
    private Header header;
    private String messageId;

    public static SubmitSmResp fromBytes(byte[] data) {
        if (data.length < Header.LENGTH) {
            throw new IllegalArgumentException("data too short to contain a submit_sm_resp PDU");
        }
        byte[] headerBytes = Arrays.copyOfRange(data, 0, Header.LENGTH);
        Header header = Header.fromBytes(headerBytes);
        String messageId = Utils.fromCOctetString(Arrays.copyOfRange(data, Header.LENGTH, data.length));
        return new SubmitSmResp(header, messageId);
    }
}
