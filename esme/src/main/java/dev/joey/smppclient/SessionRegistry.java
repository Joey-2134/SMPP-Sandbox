package dev.joey.smppclient;

import org.springframework.stereotype.Service;
import org.springframework.web.servlet.mvc.method.annotation.SseEmitter;

import java.io.IOException;
import java.util.Collection;
import java.util.UUID;
import java.util.concurrent.ConcurrentHashMap;

@Service
public class SessionRegistry {

    private final ConcurrentHashMap<UUID, SmppClient> sessions = new ConcurrentHashMap<>();
    private final ConcurrentHashMap<UUID, SseEmitter> emitters = new ConcurrentHashMap<>();

    public void registerEmitter(UUID id, SseEmitter emitter) {
        emitters.put(id, emitter);
        SmppClient client = sessions.get(id);
        if (client != null) {
            wireEventListener(id, client);
        }
    }

    private void wireEventListener(UUID id, SmppClient client) {
        client.setEventListener(event -> {
            SseEmitter emitter = emitters.get(id);
            if (emitter == null) return;
            try {
                emitter.send(event);
            } catch (IOException e) {
                emitters.remove(id);
            }
        });
    }

    public UUID createSession(String host, int port, String systemId, String password, BindType bindType) {
        SmppClient client = new SmppClient(host, port);
        try {
            client.connect();
            client.bind(systemId, password, bindType);
            sessions.put(client.getClientId(), client);
            SseEmitter emitter = emitters.get(client.getClientId());
            if (emitter != null) {
                wireEventListener(client.getClientId(), client);
            }
            return client.getClientId();
        } catch (IOException e) {
            throw new RuntimeException("Failed to create session", e);
        }
    }

    public void removeSession(UUID clientId) {
        SmppClient client = sessions.remove(clientId);
        if (client != null) {
            try {
                client.unbind();
                emitters.remove(clientId);
            } catch (IOException e) {
                throw new RuntimeException("Failed to remove session", e);
            }
        }
    }

    public SmppClient getSession(UUID clientId) {
        return sessions.get(clientId);
    }

    public Collection<SmppClient> listSessions() {
        return sessions.values();
    }

}
