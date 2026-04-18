package dev.joey.smppclient.pdu;

public class DeliverSmResp {
    private final int sequenceNumber;

    public DeliverSmResp(int sequenceNumber) {
        this.sequenceNumber = sequenceNumber;
    }

    public byte[] toBytes() {
        Header header = new Header(Header.LENGTH, CommandId.DELIVER_SM_RESP, CommandStatus.ESME_ROK, sequenceNumber);
        return header.toBytes();
    }
}
