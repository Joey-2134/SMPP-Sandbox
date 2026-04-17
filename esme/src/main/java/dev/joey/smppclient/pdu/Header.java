package dev.joey.smppclient.pdu;

import java.nio.ByteBuffer;
import java.nio.ByteOrder;

import lombok.AllArgsConstructor;
import lombok.Data;
import lombok.NoArgsConstructor;

@Data
@AllArgsConstructor
@NoArgsConstructor
public class Header {
    public static final int LENGTH = 16;
    
    private int commandLength;
    private int commandId;
    private int commandStatus;
    private int sequenceNumber;

    public byte[] toBytes() {
        return ByteBuffer.allocate(LENGTH)
            .order(ByteOrder.BIG_ENDIAN)
            .putInt(commandLength)
            .putInt(commandId)
            .putInt(commandStatus)
            .putInt(sequenceNumber)
            .array();
    }

    public static Header fromBytes(byte[] data) {
        ByteBuffer buffer = ByteBuffer.wrap(data).order(ByteOrder.BIG_ENDIAN);
        return new Header(
            buffer.getInt(),
            buffer.getInt(),
            buffer.getInt(),
            buffer.getInt()
        );
    }
}