package dev.joey.smppclient.controller;

import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;
import org.springframework.web.servlet.mvc.method.annotation.SseEmitter;
import org.springframework.web.bind.annotation.DeleteMapping;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;

import java.util.Collection;
import java.util.UUID;

import dev.joey.smppclient.SmppClient;
import dev.joey.smppclient.SessionRegistry;

@RestController
@RequestMapping("/api/sessions")
public class SessionController {
    private final SessionRegistry sessionRegistry;

    public SessionController(SessionRegistry sessionRegistry) {
        this.sessionRegistry = sessionRegistry;
    }

    @GetMapping("/{id}/events/stream")
    public SseEmitter stream(@PathVariable UUID id) {
        SseEmitter emitter = new SseEmitter(0L);
        sessionRegistry.registerEmitter(id, emitter);
        return emitter;
  }

    @GetMapping
    public Collection<SmppClient> listSessions() {
        return sessionRegistry.listSessions();
    }

    @GetMapping("/{clientId}")
    public SmppClient getSession(@PathVariable UUID clientId) {
        return sessionRegistry.getSession(clientId);
    }

    @PostMapping
    public UUID createSession(@RequestBody SessionCreateRequest request) {
        return sessionRegistry.createSession(
            request.getHost(),
            request.getPort(),
            request.getSystemId(),
            request.getPassword()
        );
    }

    @DeleteMapping("/{clientId}")
    public void deleteSession(@PathVariable UUID clientId) {
        sessionRegistry.removeSession(clientId);
    }
}
