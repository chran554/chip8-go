package se.cha.chip8.screen;

import com.fasterxml.jackson.annotation.JsonProperty;
import lombok.Data;

@Data
public class PeripheralState {
    @JsonProperty("sound")
    boolean sound;
    @JsonProperty("keys")
    int keys;
    @JsonProperty("screen")
    byte[] screen;
    @JsonProperty("screenWidth")
    int screenWidth;
    @JsonProperty("screenHeight")
    int screenHeight;
}
