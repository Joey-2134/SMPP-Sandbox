package dev.joey.smppclient.pdu;

import java.nio.charset.StandardCharsets;

public class Utils {
    public static final byte INTERFACE_VERSION = 0x34;
    
    public static byte[] toCOctetString(String str) {
        return str.getBytes(StandardCharsets.UTF_8);
    }

    public static String fromCOctetString(byte[] bytes) {
        return new String(bytes, StandardCharsets.UTF_8);
    }

}