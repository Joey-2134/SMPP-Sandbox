package dev.joey.smppclient.pdu;

public class Unbind {
    private Header header;

    public Unbind(int sequenceNumber) {
        this.header = new Header(Header.LENGTH, CommandId.UNBIND, 0, sequenceNumber);
    }

    public byte[] toBytes() {
        return header.toBytes();
    }
}
