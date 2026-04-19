package dev.joey.smppclient.pdu;

public class DeliverSmResp {
    private final int sequenceNumber;

    public DeliverSmResp(int sequenceNumber) {
        this.sequenceNumber = sequenceNumber;
    }

    public byte[] toBytes() {
        Header header = new Header(Header.LENGTH + 1, CommandId.DELIVER_SM_RESP, CommandStatus.ESME_ROK, sequenceNumber);
        byte[] bytes = new byte[Header.LENGTH + 1];
        System.arraycopy(header.toBytes(), 0, bytes, 0, Header.LENGTH);
        bytes[Header.LENGTH] = 0x00; // empty message_id, null terminated
        return bytes;
    }
}
