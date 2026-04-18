package dev.joey.smppclient.pdu;

public class EnquireLink {
    private Header header;

    public EnquireLink(int sequenceNumber) {
        this.header = new Header(Header.LENGTH, CommandId.ENQUIRE_LINK, 0, sequenceNumber);
    }

    public byte[] toBytes() {
        return header.toBytes();
    }
}
