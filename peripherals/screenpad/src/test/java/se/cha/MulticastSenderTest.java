package se.cha;

import org.junit.Test;

import java.io.IOException;
import java.net.DatagramPacket;
import java.net.DatagramSocket;
import java.net.InetAddress;

public class MulticastSenderTest {

    @Test
    public void testSend() throws InterruptedException {

        while (true) {
            final int keyState = (int) Math.round(Math.random() * 15);

            DatagramSocket socket;
            InetAddress group;
            final byte[] buffer = new byte[2];

            buffer[0] = (byte) ((keyState & 0xFF00) >>> 8);
            buffer[1] = (byte) ((keyState & 0x00FF) >>> 0);

            System.out.println("Sending key state: " +
                    leftPad(Integer.toBinaryString(buffer[0] & 0x000000FF), "0", 8) + " " +
                    leftPad(Integer.toBinaryString(buffer[1] & 0x000000FF), "0", 8) + " " +
                    leftPad(Integer.toBinaryString(keyState), "0", 16));

            try {
                socket = new DatagramSocket();
                group = InetAddress.getByName("230.0.0.0");

                final DatagramPacket packet = new DatagramPacket(buffer, buffer.length, group, 9998);
                socket.send(packet);
                socket.close();
            } catch (IOException e) {
                throw new RuntimeException(e);
            }

            Thread.sleep(3 * 1000);
        }
    }

    public static String leftPad(String text, String padding, int length) {
        while (text.length() < length) {
            text = padding + text;
        }

        return text;
    }

}
