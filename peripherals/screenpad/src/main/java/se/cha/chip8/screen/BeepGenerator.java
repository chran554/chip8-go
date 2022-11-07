package se.cha.chip8.screen;

import javax.sound.sampled.*;

public class BeepGenerator extends Thread {

    private static final int FREQUENCY = 44100;
    private static final byte[] SAMPLE_BUFFER = new byte[2]; // A two byte buffer, as we flashingly will add a harmonic (one octave higher) to the sine sound wave
    private static final AudioFormat AUDIO_FORMAT = new AudioFormat(FREQUENCY, 8, 2, true, false);

    private static SourceDataLine sourceDataLine;
    private static BeepGenerator beepGenerator;

    private boolean stop = false;

    private final int hz;
    private final double volume;


    private BeepGenerator(int hz, double volume) {
        super();

        this.hz = hz;
        this.volume = volume;

        try {
            sourceDataLine = AudioSystem.getSourceDataLine(AUDIO_FORMAT);
            sourceDataLine.open(AUDIO_FORMAT);
        } catch (LineUnavailableException e) {
            throw new RuntimeException(e);
        }
    }

    private static void createBeepGenerator() {
        if (beepGenerator == null) {
            beepGenerator = new BeepGenerator(220, 0.5);
        }
    }

    public synchronized static void playBeep() {
        if (beepGenerator == null) {
            startBeepGenerator();
        }

        sourceDataLine.start();
    }

    public synchronized static void pauseBeep() {
        sourceDataLine.stop();
    }

    public synchronized static void stopBeepGenerator() {
        // Stop audio beep. Clear that source data line and end it.
        sourceDataLine.stop();

        if (beepGenerator != null) {
            beepGenerator.stop = true;
            beepGenerator.interrupt();
            beepGenerator = null;
        }

        sourceDataLine.flush();
        sourceDataLine.close();
    }

    public synchronized static void startBeepGenerator() {
        if (beepGenerator == null) {
            createBeepGenerator();
            beepGenerator.stop = false;
            beepGenerator.start();
        }
    }

    @Override
    public void run() {
        int soundWaveIndex = 0;

        while (!stop) {
            try {
                Thread.sleep(100);
            } catch (InterruptedException e) {
                // Not interested...
            }

            final int sourceDataLineDataCapacity = sourceDataLine.available();

            // Fill the audio source data line with more delicate beeping sine-waveform sound
            for (int i = 0; i < (sourceDataLineDataCapacity / 2); i++) {
                final double angle = 1.0 * soundWaveIndex / (1.0 * FREQUENCY / hz) * (2.0 * Math.PI);
                SAMPLE_BUFFER[0] = (byte) (Math.sin(angle) * 255.0 * volume * 1.0);
                SAMPLE_BUFFER[1] = (byte) (Math.sin(2.0 * angle) * 255.0 * volume * 0.6); // Harmonics at double frequency, and lower volume

                sourceDataLine.write(SAMPLE_BUFFER, 0, SAMPLE_BUFFER.length);

                soundWaveIndex = (soundWaveIndex + 1) % FREQUENCY;
            }
        }

    }
}
