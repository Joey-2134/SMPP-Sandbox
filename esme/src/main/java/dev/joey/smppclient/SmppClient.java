package dev.joey.smppclient;

import dev.joey.smppclient.pdu.*;
import lombok.Getter;
import lombok.Setter;

import java.io.IOException;
import java.io.InputStream;
import java.io.OutputStream;
import java.net.Socket;
import java.nio.ByteBuffer;
import java.nio.ByteOrder;
import java.time.LocalDateTime;
import java.util.Arrays;
import java.util.UUID;
import java.util.concurrent.CompletableFuture;
import java.util.concurrent.ConcurrentHashMap;
import java.util.concurrent.ExecutionException;
import java.util.concurrent.TimeUnit;
import java.util.concurrent.TimeoutException;
import java.util.concurrent.atomic.AtomicInteger;
import java.util.function.Consumer;

public class SmppClient {
    @Getter
    private final UUID clientId;
    private final String host;
    private final int port;
    private Socket socket;
    private InputStream in;
    private OutputStream out;
    @Setter
    private Consumer<SessionEvent> eventListener = e -> {}; // no-op default
    private final AtomicInteger sequenceNumber = new AtomicInteger(1);
    private final ConcurrentHashMap<Integer, CompletableFuture<byte[]>> pendingResponses = new ConcurrentHashMap<>();
    private final ConcurrentHashMap<Integer, Consumer<SubmitSmResp>> submitCallbacks = new ConcurrentHashMap<>();
    

    public SmppClient(String host, int port) {
        this.clientId = UUID.randomUUID();
        this.host = host;
        this.port = port;
    }

    public void connect() throws IOException {
        socket = new Socket(host, port);
        in = socket.getInputStream();
        out = socket.getOutputStream();
        Thread readThread = new Thread(this::readLoop);
        readThread.start();
        System.out.println("Connected to " + host + ":" + port);
    }

    //blocking
    public void bind(String systemId, String password) throws IOException {
        int seqNum = nextSequenceNumber();
        CompletableFuture<byte[]> future = new CompletableFuture<>();
        pendingResponses.put(seqNum, future);

        out.write(new BindTransmitter(seqNum, systemId, password, "").toBytes());

        try {
            byte[] respBytes = future.get(5, TimeUnit.SECONDS);
            BindTransmitterResp resp = BindTransmitterResp.fromBytes(respBytes);
            if (resp.getHeader().getCommandStatus() != CommandStatus.ESME_ROK) {
                throw new IOException("Bind failed with status: 0x" + Integer.toHexString(resp.getHeader().getCommandStatus()));
            }
            System.out.println("Bound as: " + resp.getSystemId());
        } catch (InterruptedException | ExecutionException | TimeoutException e) {
            throw new IOException("Bind failed", e);
        }
    }

    //blocking
    public void unbind() throws IOException {
        int seqNum = nextSequenceNumber();
        CompletableFuture<byte[]> future = new CompletableFuture<>();
        pendingResponses.put(seqNum, future);

        out.write(new Unbind(seqNum).toBytes());

        try {
            byte[] respBytes = future.get(5, TimeUnit.SECONDS);
            UnbindResp resp = UnbindResp.fromBytes(respBytes);
            if (resp.getHeader().getCommandStatus() != CommandStatus.ESME_ROK) {
                throw new IOException("Unbind failed with status: 0x" + Integer.toHexString(resp.getHeader().getCommandStatus()));
            }
            System.out.println("Unbound");
            socket.close();
        } catch (InterruptedException | ExecutionException | TimeoutException e) {
            throw new IOException("Unbind failed", e);
        }
    }

    // fire and forget
    public void submitSm(String from, String to, String message, Consumer<SubmitSmResp> callback) throws IOException {
        int seqNum = nextSequenceNumber();
        submitCallbacks.put(seqNum, callback);
        out.write(SubmitSm.basic(seqNum, from, to, message).toBytes());
        eventListener.accept(new SessionEvent(
            SessionEvent.EventType.SUBMIT_SENT,
             LocalDateTime.now().toString(),
              "Sent to " + to + ": " + message
            ));
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

    private void readLoop() {
        try {
            while (!socket.isClosed()) {
                byte[] raw = readPdu();
                Header header = Header.fromBytes(Arrays.copyOfRange(raw, 0, Header.LENGTH));
  
                if ((header.getCommandId() & 0x80000000) != 0) {
                    // response PDU
                    CompletableFuture<byte[]> future = pendingResponses.remove(header.getSequenceNumber());
                    if (future != null) {
                        future.complete(raw);
                        continue;
                    }
                    Consumer<SubmitSmResp> callback = submitCallbacks.remove(header.getSequenceNumber());
                    if (callback != null) {
                        SubmitSmResp resp = SubmitSmResp.fromBytes(raw);
                        eventListener.accept(new SessionEvent(
                            SessionEvent.EventType.SUBMIT_ACKED,
                            LocalDateTime.now().toString(),
                            "Message ID: " + resp.getMessageId()
                        ));
                        callback.accept(resp);
                    }
                } else {
                    switch (header.getCommandId()) {
                        case CommandId.DELIVER_SM -> handleDeliverSm(header, raw);
                        case CommandId.ENQUIRE_LINK -> handleEnquireLink(header);
                        default -> handleDefault(header);
                    }
                }
            }
        } catch (IOException e) {
            pendingResponses.values().forEach(f -> f.completeExceptionally(e));
            submitCallbacks.clear();
        }
    }
  
    private void handleDeliverSm(Header header, byte[] raw) throws IOException {
        SessionEvent event = new SessionEvent(
            SessionEvent.EventType.DELIVER_SM,
            LocalDateTime.now().toString(),
            "Received deliver_sm PDU"
        );
        eventListener.accept(event);
        System.out.println("Received deliver_sm PDU");
        out.write(new DeliverSmResp(header.getSequenceNumber()).toBytes());
    }
  
    private void handleEnquireLink(Header header) throws IOException {
        out.write(new EnquireLinkResp(header).toBytes());
    }
  
    private void handleDefault(Header header) throws IOException {
        out.write(new GenericNack(header.getSequenceNumber(), CommandStatus.ESME_RINVCMDID).toBytes());
    }
}
