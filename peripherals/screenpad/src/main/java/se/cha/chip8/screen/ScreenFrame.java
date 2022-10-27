package se.cha.chip8.screen;

import javax.swing.*;
import java.awt.*;
import java.awt.image.BufferedImage;

public class ScreenFrame extends JFrame {

    private static ScreenFrame singleton = null;

    protected BufferedImage image = null;
    protected JScrollPane scrollPane = new JScrollPane(JScrollPane.VERTICAL_SCROLLBAR_AS_NEEDED, JScrollPane.HORIZONTAL_SCROLLBAR_AS_NEEDED);

    private int[] bufferImageData;
    private BufferedImage bufferImage;

    private int width = 64;
    private int height = 32;

    private ScreenFrame() {
        super("CHIP-8");
        initUI();
        setDefaultCloseOperation(JFrame.EXIT_ON_CLOSE);

        initialize(64, 32);
        pack();
    }

    public static ScreenFrame getOrCreateSingleton() {
        if (singleton == null) {
            singleton = new ScreenFrame();
            singleton.setVisible(true);
            singleton.centerFrame();
        }

        return singleton;
    }

    protected void initUI() {
        add(scrollPane);
    }

    public BufferedImage getImage() {
        return image;
    }

    public void setImage(BufferedImage image) {
        this.image = image;

        final ImageIcon imageIcon = new ImageIcon(image);
        final JLabel label = new JLabel(imageIcon);
        scrollPane.setViewportView(label);

        resizeFrame();
    }


    public synchronized void setChip8ScreenData(byte[] imageBitData) {
        // Convert bit array (of bytes) with one bit per pixel to array of int with one int per pixel
        for (int y = 0; y < height; y++) {
            for (int x = 0; x < width; x++) {
                final int streamPixelIndex = x + (y * width);
                final int streamByteIndex = streamPixelIndex / 8;
                final int streamBitInByteIndex = streamPixelIndex % 8;

                final int streamBitValue = (imageBitData[streamByteIndex] >> (7 - streamBitInByteIndex)) & 0x00000001;
                bufferImageData[streamPixelIndex] = streamBitValue == 1 ? 0xFFFFFFFF : 0xFF000000;
            }
        }

        bufferImage.setRGB(0, 0, width, height, bufferImageData, 0, width);
        final Graphics2D g = (Graphics2D) getImage().getGraphics();
        g.drawImage(bufferImage, 0, 0, width * 10, height * 10, null);

        repaint();
    }

    public void initialize(int width, int height) {
        setImage(new BufferedImage(width * 10, height * 10, BufferedImage.TYPE_BYTE_BINARY));

        bufferImage = new BufferedImage(width, height, BufferedImage.TYPE_BYTE_BINARY);
        bufferImageData = new int[width * height];

        final Graphics graphics = image.getGraphics();
        graphics.setColor(java.awt.Color.BLACK);
        graphics.fillRect(0, 0, width, height);
    }

    protected void resizeFrame() {
        // Make sure window fit within desktop area
        final Rectangle maxWindowBounds = GraphicsEnvironment.getLocalGraphicsEnvironment().getMaximumWindowBounds();
        final int width = getWidth();
        final int height = getHeight();
        setSize(Math.min(width, maxWindowBounds.width), Math.min(height, maxWindowBounds.height));
    }

    public void centerFrame() {
        setLocationRelativeTo(null);
    }
}