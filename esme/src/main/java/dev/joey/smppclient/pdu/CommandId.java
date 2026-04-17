package dev.joey.smppclient.pdu;

public class CommandId {
    public static final int GENERIC_NACK = 0x80000000;
    public static final int BIND_RECEIVER = 0x00000001;
    public static final int BIND_RECEIVER_RESP = 0x80000001;
    public static final int BIND_TRANSMITTER = 0x00000002;
    public static final int BIND_TRANSMITTER_RESP = 0x80000002;
    public static final int BIND_TRANSCEIVER = 0x00000009;
    public static final int BIND_TRANSCEIVER_RESP = 0x80000009;
    public static final int SUBMIT_SM = 0x00000004;
    public static final int SUBMIT_SM_RESP = 0x80000004;
    public static final int DELIVER_SM = 0x00000005;
    public static final int DELIVER_SM_RESP = 0x80000005;
    public static final int ENQUIRE_LINK = 0x00000015;
    public static final int ENQUIRE_LINK_RESP = 0x80000015;
    public static final int UNBIND = 0x00000006;
    public static final int UNBIND_RESP = 0x80000006;
}