package dev.joey.smppclient;

import lombok.AllArgsConstructor;
import lombok.Getter;

@Getter
@AllArgsConstructor
public class SessionEvent {
    public enum EventType {
        SUBMIT_SENT,
        SUBMIT_ACKED,
        DELIVER_SM,
        ERROR
    }

    private EventType type;
    private String timestamp;
    private String message;
}
