package dev.joey.smppclient;

import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;

@SpringBootApplication
public class Main {
    public static void main(String[] args) throws Exception {
        SpringApplication.run(Main.class, args);
        System.out.println("SMPP client started");
        System.out.println("Frontend is available at http://localhost:8080");
        System.in.read();
        System.exit(0);
    }
}
