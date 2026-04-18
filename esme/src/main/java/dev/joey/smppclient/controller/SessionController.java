package dev.joey.smppclient.controller;

import org.springframework.http.HttpStatus;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.ResponseStatus;
import org.springframework.web.bind.annotation.RestController;
import org.springframework.web.server.ResponseStatusException;
import org.springframework.web.servlet.mvc.method.annotation.SseEmitter;
import org.springframework.web.bind.annotation.DeleteMapping;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;

import java.io.IOException;
import java.util.Collection;
import java.util.UUID;

import dev.joey.smppclient.BindType;
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
        BindType bindType = request.getBindType() != null ? request.getBindType() : BindType.TRX;
        return sessionRegistry.createSession(
            request.getHost(),
            request.getPort(),
            request.getSystemId(),
            request.getPassword(),
            bindType
        );
    }

    @DeleteMapping("/{clientId}")
    public void deleteSession(@PathVariable UUID clientId) {
        sessionRegistry.removeSession(clientId);
    }

    @PostMapping("/{id}/submit")
    @ResponseStatus(HttpStatus.ACCEPTED)
    public void submit(@PathVariable UUID id, @RequestBody SubmitRequest request) {
        SmppClient client = sessionRegistry.getSession(id);
        if (client == null) {
            throw new ResponseStatusException(HttpStatus.NOT_FOUND, "Session not found");
        }
        if (client.getBindType() == BindType.RX) {
            throw new ResponseStatusException(HttpStatus.METHOD_NOT_ALLOWED, "RX session cannot send messages");
        }
        try {
            client.submitSm(request.getFrom(), request.getTo(), request.getMessage(), resp -> {});
        } catch (IOException e) {
            throw new ResponseStatusException(HttpStatus.INTERNAL_SERVER_ERROR, "Submit failed: " + e.getMessage());
        }
    }
}
