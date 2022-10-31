package se.cha.chip8.screen;

import java.awt.*;
import java.awt.image.BufferedImage;
import java.awt.image.ConvolveOp;
import java.awt.image.Kernel;

public class Blur2 {

    public static BufferedImage blur(BufferedImage image, int radius, boolean fasterBlur) {
        if (fasterBlur) {
            image = changeImageWidth(image, image.getWidth() / 2);
            image = getGaussianBlurFilter(radius / 2, true).filter(image, null);
            image = getGaussianBlurFilter(radius / 2, false).filter(image, null);
            image = changeImageWidth(image, image.getWidth() * 2);
        } else {
            image = getGaussianBlurFilter(radius, true).filter(image, null);
            image = getGaussianBlurFilter(radius, false).filter(image, null);
        }

        return image;
    }

    public static ConvolveOp getGaussianBlurFilter(int radius, boolean horizontal) {
        if (radius < 1) {
            throw new IllegalArgumentException("Radius must be >= 1");
        }

        final int size = radius * 2 + 1;
        final float[] data = new float[size];

        final float sigma = radius / 3.0f;
        final float twoSigmaSquare = 2.0f * sigma * sigma;
        final float sigmaRoot = (float) Math.sqrt(twoSigmaSquare * Math.PI);
        float total = 0.0f;

        for (int i = -radius; i <= radius; i++) {
            final float distance = i * i;
            final int index = i + radius;
            data[index] = (float) Math.exp(-distance / twoSigmaSquare) / sigmaRoot;
            total += data[index];
        }

        for (int i = 0; i < data.length; i++) {
            data[i] /= total;
        }

        final Kernel kernel = horizontal ? new Kernel(size, 1, data) : new Kernel(1, size, data);

        return new ConvolveOp(kernel, ConvolveOp.EDGE_ZERO_FILL, null);
    }

    public static BufferedImage changeImageWidth(BufferedImage image, int width) {
        final float ratio = (float) image.getWidth() / (float) image.getHeight();
        final int height = (int) (width / ratio);

        final BufferedImage temp = new BufferedImage(width, height, image.getType());
        final Graphics2D g2 = temp.createGraphics();
        g2.setRenderingHint(RenderingHints.KEY_INTERPOLATION, RenderingHints.VALUE_INTERPOLATION_BILINEAR);
        g2.drawImage(image, 0, 0, temp.getWidth(), temp.getHeight(), null);
        g2.dispose();

        return temp;
    }
}
