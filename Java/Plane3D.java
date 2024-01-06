package Project;

public class Plane3D {

    private double a, b, c, d;

    public Plane3D(double a, double b, double c, double d) {

        this.a = a;
        this.b = b;
        this.c = c;
        this.d = d;
    }

    public Plane3D(Point3D p1, Point3D p2, Point3D p3) {

        this.a = (p2.getY() - p1.getY()) * (p3.getZ() - p1.getZ()) - (p3.getY() - p1.getY()) * (p2.getZ() - p1.getZ());
        this.b = (p2.getZ() - p1.getZ()) * (p3.getX() - p1.getX()) - (p3.getZ() - p1.getZ()) * (p2.getX() - p1.getX());
        this.c = (p2.getX() - p1.getX()) * (p3.getY() - p1.getY()) - (p3.getX() - p1.getX()) * (p2.getY() - p1.getY());
        this.d = -(a * p1.getX() + b * p1.getY() + c * p1.getZ());
    }

    public double distanceFromPoint(Point3D point) {

        double num = a * point.getX() + b * point.getY() + c * point.getZ() + d;
        double denom = Math.sqrt(a * a + b * b + c * c);
        return num / denom;
    }

    public double getA() { return a; }
    public double getB() { return b; }
    public double getC() { return c; }
    public double getD() { return d; }

}
