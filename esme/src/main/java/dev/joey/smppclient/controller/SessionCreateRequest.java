package dev.joey.smppclient.controller;

import lombok.Data;
import lombok.NoArgsConstructor;

@Data
@NoArgsConstructor
public class SessionCreateRequest {
    private String host;
    private int port;
    private String systemId;
    private String password;
}
