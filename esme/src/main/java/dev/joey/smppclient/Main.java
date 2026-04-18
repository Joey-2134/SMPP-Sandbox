package dev.joey.smppclient;

import java.io.IOException;

public class Main {
    public static void main(String[] args) {
        SmppClient client = new SmppClient("localhost", 2775);
        try {
            client.connect();
            client.bind("joeysSystemId", "password");
            client.submitSm("joey", "smsc", "Hello, World!");
            client.unbind();
        } catch (IOException e) {
            e.printStackTrace();
        }
    }
}