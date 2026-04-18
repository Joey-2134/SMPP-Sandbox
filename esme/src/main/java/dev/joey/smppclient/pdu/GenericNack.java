package dev.joey.smppclient.pdu;

public class GenericNack {
    private final int sequenceNumber;
    private final int commandStatus;

    public GenericNack(int sequenceNumber, int commandStatus) {
        this.sequenceNumber = sequenceNumber;
        this.commandStatus = commandStatus;
    }

    public byte[] toBytes() {
        Header header = new Header(Header.LENGTH, CommandId.GENERIC_NACK, commandStatus, sequenceNumber);
        return header.toBytes();
    }
}
