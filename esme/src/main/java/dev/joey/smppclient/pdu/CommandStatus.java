package dev.joey.smppclient.pdu;

public class CommandStatus {
    public static final int ESME_ROK = 0x00000000; // Success
    public static final int ESME_RINVCMDID = 0x00000003; // Invalid command ID
    public static final int ESME_RINVBNDSTS = 0x00000004; // Invalid bind status
}