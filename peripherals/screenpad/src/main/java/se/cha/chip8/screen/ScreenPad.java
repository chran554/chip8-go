package se.cha.chip8.screen;

public class ScreenPad {

    public static void main(String[] args) {
        System.out.println("Listening CHIP-8 screen broadcast...");
        final UDPDataProcessor renderMessageProcessor = new UDPDataProcessor();
        final UDPMulticastMessageListener dataListener = new UDPMulticastMessageListener(renderMessageProcessor, "230.0.0.0", 9999);

        final Thread messageThread = new Thread(dataListener);
        messageThread.start();
    }

}
