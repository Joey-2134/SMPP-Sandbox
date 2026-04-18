package dev.joey.smppclient.pdu;

import java.io.ByteArrayOutputStream;
import java.io.IOException;

import lombok.AllArgsConstructor;
import lombok.Data;
import lombok.NoArgsConstructor;

@Data
@AllArgsConstructor
@NoArgsConstructor
public class SubmitSm {
    private Header header;
    private String serviceType;
    private int sourceAddrTon;
    private int sourceAddrNpi;
    private String sourceAddr;
    private int destAddrTon;
    private int destAddrNpi;
    private String destAddr;
    private int esmClass;
    private int protocolId;
    private int priorityFlag;
    private String scheduleDeliveryTime;
    private String validityPeriod;
    private int registeredDelivery;
    private int replaceIfPresentFlag;
    private int dataCoding;
    private int smDefaultMsgId;
    private byte[] message;

    public byte[] toBytes() { 
        byte[] serviceTypeBytes = Utils.toCOctetString(serviceType);
        byte[] sourceAddrBytes = Utils.toCOctetString(sourceAddr);
        byte[] destAddrBytes = Utils.toCOctetString(destAddr);
        byte[] scheduleDeliveryTimeBytes = Utils.toCOctetString(scheduleDeliveryTime);
        byte[] validityPeriodBytes = Utils.toCOctetString(validityPeriod); 
      
        int bodyLength = serviceTypeBytes.length + sourceAddrBytes.length +
            destAddrBytes.length+ scheduleDeliveryTimeBytes.length +
            validityPeriodBytes.length + message.length + 12;
      
        Header header = new Header(
            Header.LENGTH + bodyLength,
            CommandId.SUBMIT_SM,
            0,
            this.header.getSequenceNumber()
        );
      
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
