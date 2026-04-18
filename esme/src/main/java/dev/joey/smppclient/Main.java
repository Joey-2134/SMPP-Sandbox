package dev.joey.smppclient;

import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;

@SpringBootApplication
public class Main {
    public static void main(String[] args) throws Exception {
        SpringApplication.run(Main.class, args);

        SmppClient client = new SmppClient("localhost", 2775);
        client.connect();
        client.bind("joeysSystemId", "password");
        client.submitSm("joey", "smsc", "Hello, World!", resp ->
                System.out.println("Submit ack received, message ID: " + resp.getMessageId())
        );
        Thread.sleep(2000);
        client.unbind();
    }
}
