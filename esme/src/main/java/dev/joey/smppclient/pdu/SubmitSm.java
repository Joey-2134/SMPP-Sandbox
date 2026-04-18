package dev.joey.smppclient.pdu;

import java.io.ByteArrayOutputStream;
import java.io.IOException;
import java.nio.charset.StandardCharsets;

import lombok.AllArgsConstructor;

@AllArgsConstructor
public class SubmitSm {
    private final int sequenceNumber;
    private final String serviceType;
    private final int sourceAddrTon;
    private final int sourceAddrNpi;
    private final String sourceAddr;
    private final int destAddrTon;
    private final int destAddrNpi;
    private final String destAddr;
    private final int esmClass;
    private final int protocolId;
    private final int priorityFlag;
    private final String scheduleDeliveryTime;
    private final String validityPeriod;
    private final int registeredDelivery;
    private final int replaceIfPresentFlag;
    private final int dataCoding;
    private final int smDefaultMsgId;
    private final byte[] message;

    public static SubmitSm basic(int sequenceNumber, String sourceAddr, String destAddr, String message) {
        return new SubmitSm(
                sequenceNumber,
                "",     // service_type
                0x00,   // source_addr_ton
                0x00,   // source_addr_npi
                sourceAddr,
                0x01,   // dest_addr_ton (international)
                0x01,   // dest_addr_npi (ISDN)
                destAddr,
                0x00,   // esm_class
                0x00,   // protocol_id
                0x00,   // priority_flag
                "",     // schedule_delivery_time (immediate)
                "",     // validity_period (default)
                0x00,   // registered_delivery
                0x00,   // replace_if_present_flag
                0x00,   // data_coding (GSM7)
                0x00,   // sm_default_msg_id
                message.getBytes(StandardCharsets.UTF_8)
        );
    }

    public byte[] toBytes() {
        byte[] serviceTypeBytes = Utils.toCOctetString(serviceType);
        byte[] sourceAddrBytes = Utils.toCOctetString(sourceAddr);
        byte[] destAddrBytes = Utils.toCOctetString(destAddr);
        byte[] scheduleDeliveryTimeBytes = Utils.toCOctetString(scheduleDeliveryTime);
        byte[] validityPeriodBytes = Utils.toCOctetString(validityPeriod);

        int bodyLength = serviceTypeBytes.length + sourceAddrBytes.length +
                destAddrBytes.length + scheduleDeliveryTimeBytes.length +
                validityPeriodBytes.length + message.length + 12;

        Header header = new Header(Header.LENGTH + bodyLength, CommandId.SUBMIT_SM, 0, sequenceNumber);

        ByteArrayOutputStream out = new ByteArrayOutputStream();
        try {
            out.write(header.toBytes());
            out.write(serviceTypeBytes);
            out.write(sourceAddrTon);
            out.write(sourceAddrNpi);
            out.write(sourceAddrBytes);
            out.write(destAddrTon);
            out.write(destAddrNpi);
            out.write(destAddrBytes);
            out.write(esmClass);
            out.write(protocolId);
            out.write(priorityFlag);
            out.write(scheduleDeliveryTimeBytes);
            out.write(validityPeriodBytes);
            out.write(registeredDelivery);
            out.write(replaceIfPresentFlag);
            out.write(dataCoding);
            out.write(smDefaultMsgId);
            out.write(message.length);
            out.write(message);
        } catch (IOException e) {
            throw new RuntimeException(e);
        }
        return out.toByteArray();
    }
}
