package se.cha;

/**
 * Rigorous Test :-)
 */

import org.junit.Test;

import javax.swing.*;
import java.awt.*;
import java.awt.image.BufferedImage;

public class ImageAlphaTest {

    @Test
    public void shouldAnswerWithTrue() {

        final BufferedImage imageA = new BufferedImage(300, 300, BufferedImage.TYPE_INT_ARGB);
        final BufferedImage imageB = new BufferedImage(300, 300, BufferedImage.TYPE_INT_ARGB);

        final Graphics graphicsA = imageA.getGraphics();
        graphicsA.setColor(new Color(0xFF, 0x00, 0x00, 0xAA));
        graphicsA.fillRect(0, 0, imageA.getWidth() / 2, imageA.getHeight() / 2);

        final Graphics graphicsB = imageB.getGraphics();
        graphicsB.setColor(Color.WHITE);
        graphicsB.fillRect(0, 0, imageB.getWidth(), imageB.getHeight());
        graphicsB.setColor(new Color(0x00, 0xFF, 0x00, 0xFF));
        graphicsB.fillRect(imageB.getWidth() / 3, imageB.getHeight() / 3, imageB.getWidth() * 2 / 3, imageB.getHeight() *2 / 3);

        preMultiplyAlpha(imageA);
        preMultiplyAlpha(imageB);
        blur(imageA);
        blur(imageA);
        blur(imageA);
        final BufferedImage imageC = overlay(imageA, 0, 0, imageB);
        unpreMultiplyAlpha(imageC);

        final JFrame frame = new JFrame() {

            @Override
            public void paint(Graphics g) {
                super.paint(g);
                System.out.println("paint");
                final Graphics2D graphics2D = (Graphics2D) g;
                graphics2D.drawImage(imageC, 0, 0, null);
            }
        };
        frame.setSize(350, 350);
        frame.setVisible(true);
        frame.repaint();
        frame.setDefaultCloseOperation(JFrame.EXIT_ON_CLOSE);


        try {
            Thread.sleep(30 * 1000);
        } catch (InterruptedException e) {
            throw new RuntimeException(e);
        }
        System.out.println("test");
    }

    private void blur(BufferedImage image) {
        // TODO implement
    }

    private BufferedImage overlay(BufferedImage imageA, int x, int y, BufferedImage imageB) {
        final BufferedImage imageC = new BufferedImage(imageB.getWidth(), imageB.getWidth(), BufferedImage.TYPE_INT_ARGB);
        final int[] rgbC = imageC.getRGB(0, 0, imageC.getWidth(), imageC.getHeight(), null, 0, imageC.getWidth());
        final int[] rgbB = imageB.getRGB(0, 0, imageB.getWidth(), imageB.getHeight(), null, 0, imageB.getWidth());
        final int[] rgbA = imageA.getRGB(0, 0, imageA.getWidth(), imageA.getHeight(), null, 0, imageA.getWidth());

        for (int py = 0; py < imageB.getHeight(); py++) {
            for (int px = 0; px < imageB.getHeight(); px++) {
                if (((px < x) || (px >= imageA.getWidth() + x)) || ((py < y) || (py >= imageA.getHeight() + y))) {
                    rgbC[px + py * imageB.getWidth()] = rgbB[px + py * imageB.getWidth()];
                } else {
                    final int pixelIndexA = (px - x) + (py - y) * imageA.getWidth();
                    final int pixelIndexB = px + py * imageB.getWidth();

                    final float aa = ((rgbA[pixelIndexA] & 0xFF000000) >>> 24) / 255.0f;
                    final float ar = ((rgbA[pixelIndexA] & 0x00FF0000) >> 16) / 255.0f;
                    final float ag = ((rgbA[pixelIndexA] & 0x0000FF00) >> 8) / 255.0f;
                    final float ab = ((rgbA[pixelIndexA] & 0x000000FF) >> 0) / 255.0f;

                    final float ba = ((rgbB[pixelIndexB] & 0xFF000000) >>> 24) / 255.0f;
                    final float br = ((rgbB[pixelIndexB] & 0x00FF0000) >> 16) / 255.0f;
                    final float bg = ((rgbB[pixelIndexB] & 0x0000FF00) >> 8) / 255.0f;
                    final float bb = ((rgbB[pixelIndexB] & 0x000000FF) >> 0) / 255.0f;

                    final float k = (1.0f - aa);
                    final float ca = aa + ba * k;
                    final float cr = ar + br * k;
                    final float cg = ag + bg * k;
                    final float cb = ab + bb * k;

                    rgbC[pixelIndexB] = // rgbB[pixelIndexB];
                            (Math.round(trunc(ca, 0.0f, 1.0f) * 255) << 24) |
                            (Math.round(trunc(cr, 0.0f, 1.0f) * 255) << 16) |
                            (Math.round(trunc(cg, 0.0f, 1.0f) * 255) << 8) |
                            (Math.round(trunc(cb, 0.0f, 1.0f) * 255) << 0);
                }
            }
        }


        imageC.setRGB(0, 0, imageC.getWidth(), imageC.getHeight(), rgbC, 0, imageC.getWidth());
        return imageC;
    }

    private float trunc(float v, float min, float max) {
        return Math.max(min, Math.min(v, max));
    }

    private void unpreMultiplyAlpha(BufferedImage image) {
        final int[] rgb = image.getRGB(0, 0, image.getWidth(), image.getHeight(), null, 0, image.getWidth());

        for (int i = 0; i < rgb.length; i++) {
            final float a = ((rgb[i] & 0xFF000000) >>> 24) / 255.0f;
            final float r = ((rgb[i] & 0x00FF0000) >> 16) / 255.0f;
            final float g = ((rgb[i] & 0x0000FF00) >> 8) / 255.0f;
            final float b = ((rgb[i] & 0x000000FF) >> 0) / 255.0f;

            final float k = (a > 0.0f) ? a : 1.0f;
            rgb[i] =
                    (Math.round(a * 255) << 24) |
                            (Math.round((r / k) * 255) << 16) |
                            (Math.round((g / k) * 255) << 8) |
                            (Math.round((b / k) * 255) << 0);
        }

        image.setRGB(0, 0, image.getWidth(), image.getHeight(), rgb, 0, image.getWidth());
    }

    private void preMultiplyAlpha(BufferedImage image) {
        final int[] rgb = image.getRGB(0, 0, image.getWidth(), image.getHeight(), null, 0, image.getWidth());

        for (int i = 0; i < rgb.length; i++) {
            final float a = ((rgb[i] & 0xFF000000) >>> 24) / 255.0f;
            final float r = ((rgb[i] & 0x00FF0000) >> 16) / 255.0f;
            final float g = ((rgb[i] & 0x0000FF00) >> 8) / 255.0f;
            final float b = ((rgb[i] & 0x000000FF) >> 0) / 255.0f;

            rgb[i] =
                    (Math.round(a * 255) << 24) |
                            (Math.round((r * a) * 255) << 16) |
                            (Math.round((g * a) * 255) << 8) |
                            (Math.round((b * a) * 255) << 0);
        }

        image.setRGB(0, 0, image.getWidth(), image.getHeight(), rgb, 0, image.getWidth());
    }


}
