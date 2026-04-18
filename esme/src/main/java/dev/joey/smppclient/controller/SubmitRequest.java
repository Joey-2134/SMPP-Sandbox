package dev.joey.smppclient.controller;

import lombok.Data;
import lombok.NoArgsConstructor;

@Data
@NoArgsConstructor
public class SubmitRequest {
    private String from;
    private String to;
    private String message;
}
