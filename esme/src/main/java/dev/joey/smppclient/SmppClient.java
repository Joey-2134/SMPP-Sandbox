package dev.joey.smppclient;

import dev.joey.smppclient.pdu.*;

import java.io.IOException;
import java.io.InputStream;
import java.io.OutputStream;
import java.net.Socket;
import java.nio.ByteBuffer;
import java.nio.ByteOrder;
import java.util.concurrent.atomic.AtomicInteger;

public class SmppClient {
    private final String host;
    private final int port;
    private Socket socket;
    private InputStream in;
    private OutputStream out;
    private final AtomicInteger sequenceNumber = new AtomicInteger(1);

    public SmppClient(String host, int port) {
        this.host = host;
        this.port = port;
    }

    public void connect() throws IOException {
        socket = new Socket(host, port);
        in = socket.getInputStream();
        out = socket.getOutputStream();
        System.out.println("Connected to " + host + ":" + port);
    }

    public void bind(String systemId, String password) throws IOException {
        BindTransmitter request = new BindTransmitter(
                new Header(0, CommandId.BIND_TRANSMITTER, 0, nextSequenceNumber()),
                systemId,
                password,
                "",
                0x34,
                0x00,
                0x00,
                ""
        );
        out.write(request.toBytes());

        byte[] respBytes = readPdu();
        BindTransmitterResp resp = BindTransmitterResp.fromBytes(respBytes);
        if (resp.getHeader().getCommandStatus() != CommandStatus.ESME_ROK) {
            throw new IOException("Bind failed with status: 0x" + Integer.toHexString(resp.getHeader().getCommandStatus()));
        }
        System.out.println("Bound as: " + resp.getSystemId());
    }

    public void unbind() throws IOException {
        Unbind request = new Unbind(nextSequenceNumber());
        out.write(request.toBytes());

        byte[] respBytes = readPdu();
        UnbindResp resp = UnbindResp.fromBytes(respBytes);
        if (resp.getHeader().getCommandStatus() != CommandStatus.ESME_ROK) {
            throw new IOException("Unbind failed with status: 0x" + Integer.toHexString(resp.getHeader().getCommandStatus()));
        }
        System.out.println("Unbound");
        socket.close();
    }

    private byte[] readPdu() throws IOException {
        byte[] lenBytes = in.readNBytes(4);
        if (lenBytes.length < 4) {
            throw new IOException("Connection closed while reading PDU length");
        }
        int commandLength = ByteBuffer.wrap(lenBytes).order(ByteOrder.BIG_ENDIAN).getInt();
        byte[] rest = in.readNBytes(commandLength - 4);
        byte[] pdu = new byte[commandLength];
        System.arraycopy(lenBytes, 0, pdu, 0, 4);
        System.arraycopy(rest, 0, pdu, 4, rest.length);
        return pdu;
    }

    private int nextSequenceNumber() {
        return sequenceNumber.getAndIncrement();
    }
}
