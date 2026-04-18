package dev.joey.smppclient;

import org.springframework.stereotype.Service;

import java.io.IOException;
import java.util.Collection;
import java.util.UUID;
import java.util.concurrent.ConcurrentHashMap;

@Service
public class SessionRegistry {

    private final ConcurrentHashMap<UUID, SmppClient> sessions = new ConcurrentHashMap<>();
    

    public UUID createSession(String host, int port, String systemId, String password) {
        SmppClient client = new SmppClient(host, port);
        try {
            client.connect();
            client.bind(systemId, password);
            sessions.put(client.getClientId(), client);
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
