package dev.joey.smppclient.pdu;

import java.nio.charset.StandardCharsets;

public class Utils {
    public static final byte INTERFACE_VERSION = 0x34;
    
    public static byte[] toCOctetString(String s) {
        byte[] strBytes = s.getBytes(StandardCharsets.US_ASCII);
        byte[] result = new byte[strBytes.length + 1];
        System.arraycopy(strBytes, 0, result, 0, strBytes.length);
        result[strBytes.length] = 0x00;
        return result;
    }

    public static String fromCOctetString(byte[] bytes) {
        int len = 0;
        while (len < bytes.length && bytes[len] != 0x00) len++;
        return new String(bytes, 0, len, StandardCharsets.US_ASCII);
    }

}