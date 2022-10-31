package se.cha.chip8.screen;

import java.awt.event.KeyEvent;
import java.awt.event.KeyListener;
import java.io.IOException;
import java.net.*;
import java.util.HashMap;
import java.util.Map;

public class KeyPad implements KeyListener {

    private final Map<Integer, Integer> keyCodeToKeyPadCode = new HashMap<>();

    private int keyState = 0;

    public KeyPad() {
        keyCodeToKeyPadCode.put(KeyEvent.VK_1, 0x1);
        keyCodeToKeyPadCode.put(KeyEvent.VK_2, 0x2);
        keyCodeToKeyPadCode.put(KeyEvent.VK_3, 0x4);
        keyCodeToKeyPadCode.put(KeyEvent.VK_4, 0xC);

        keyCodeToKeyPadCode.put(KeyEvent.VK_Q, 0x4);
        keyCodeToKeyPadCode.put(KeyEvent.VK_W, 0x5);
        keyCodeToKeyPadCode.put(KeyEvent.VK_E, 0x6);
        keyCodeToKeyPadCode.put(KeyEvent.VK_R, 0xD);

        keyCodeToKeyPadCode.put(KeyEvent.VK_A, 0x7);
        keyCodeToKeyPadCode.put(KeyEvent.VK_S, 0x8);
        keyCodeToKeyPadCode.put(KeyEvent.VK_D, 0x9);
        keyCodeToKeyPadCode.put(KeyEvent.VK_F, 0xE);

        keyCodeToKeyPadCode.put(KeyEvent.VK_Z, 0xA);
        keyCodeToKeyPadCode.put(KeyEvent.VK_X, 0x0);
        keyCodeToKeyPadCode.put(KeyEvent.VK_C, 0xB);
        keyCodeToKeyPadCode.put(KeyEvent.VK_V, 0xF);

        // Chip-8 Key  Keyboard (layout QWERTY)
        // ----------  ---------
        //   1 2 3 C    1 2 3 4
        //   4 5 6 D    q w e r
        //   7 8 9 E    a s d f
        //   A 0 B F    z x c v

        // '1': 0x01, '2': 0x02, '3': 0x03, '4': 0x0C,
        // 'q': 0x04, 'w': 0x05, 'e': 0x06, 'r': 0x0D,
        // 'a': 0x07, 's': 0x08, 'd': 0x09, 'f': 0x0E,
        // 'z': 0x0A, 'x': 0x00, 'c': 0x0B, 'v': 0x0F,
    }

    public void announceKeyState() {
        DatagramSocket socket;
        InetAddress group;
        final byte[] buffer = new byte[2];

        buffer[0] = (byte) ((keyState & 0xFF00) >>> 8);
        buffer[1] = (byte) ((keyState & 0x00FF) >>> 0);

        //System.out.println("Sending key state: " +
        //        leftPad(Integer.toBinaryString(buffer[0] & 0x000000FF), "0", 8) + " " +
        //        leftPad(Integer.toBinaryString(buffer[1] & 0x000000FF), "0", 8) + " " +
        //        leftPad(Integer.toBinaryString(keyState), "0", 16));

        try {
            socket = new DatagramSocket();
            group = InetAddress.getByName("230.0.0.0");

            final DatagramPacket packet = new DatagramPacket(buffer, buffer.length, group, 9998);
            socket.send(packet);
            socket.close();
        } catch (IOException e) {
            throw new RuntimeException(e);
        }
    }

    public static String leftPad(String text, String padding, int length) {
        while (text.length() < length) {
            text = padding + text;
        }

        return text;
    }


    @Override
    public void keyTyped(KeyEvent e) {
        // Nothing by intention
    }

    @Override
    public void keyPressed(KeyEvent e) {
        final int keyCode = e.getKeyCode();

        if (keyCodeToKeyPadCode.containsKey(keyCode)) {
            final Integer padCode = keyCodeToKeyPadCode.get(keyCode);
            final int keyPadCodeBitmask = 1 << padCode;
            if ((keyState & keyPadCodeBitmask) == 0) {
                keyState |= keyPadCodeBitmask;
                //System.out.println("Pressed key: " + e.getKeyChar() + "   " + e.getKeyCode());
                announceKeyState();
            }
        }
    }

    @Override
    public void keyReleased(KeyEvent e) {
        final int keyCode = e.getKeyCode();

        if (keyCodeToKeyPadCode.containsKey(keyCode)) {
            final Integer padCode = keyCodeToKeyPadCode.get(keyCode);
            keyState &= 0xFFFF ^ (1 << padCode);
            //System.out.println("Released key: " + e.getKeyChar() + "   " + e.getExtendedKeyCode());
            announceKeyState();
        }
    }
}