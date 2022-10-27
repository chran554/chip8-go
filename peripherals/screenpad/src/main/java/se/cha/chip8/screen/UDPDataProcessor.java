package se.cha.chip8.screen;

import com.fasterxml.jackson.databind.ObjectMapper;
import org.msgpack.jackson.dataformat.MessagePackFactory;

import java.io.IOException;

public class UDPDataProcessor implements UDPMulticastMessageListener.UDPPacketDataProcessor {

    @Override
    public void onPacketReceived(byte[] data) {
        try {
            //System.out.println(new String(receivedPacket.getData(), 0, receivedPacket.getLength(), StandardCharsets.UTF_8));

            final ObjectMapper objectMapper = new ObjectMapper(new MessagePackFactory());
            final PeripheralState pixelData = objectMapper.readValue(data, PeripheralState.class);

            final ScreenFrame screenFrame = ScreenFrame.getOrCreateSingleton();

            screenFrame.setChip8ScreenData(pixelData.screen);

        } catch (IOException e) {
            e.printStackTrace();
        }
    }
}
