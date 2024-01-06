package Project;

public class Point3D {
    
    private double x,y,z;

    public Point3D(double x, double y, double z){
        this.x = x;
        this.y = y;
        this.z = z;
    }

    public boolean equals(Point3D p){
        if(p.getX() != this.x){
            return false;
        }
        if(p.getY() != this.y){
            return false;
        }
        if(p.getZ() != this.z){
            return false;
        }
        return true;
    }

    public double getX(){ return this.x; }
    public double getY(){ return this.y; }
    public double getZ(){ return this.z; }
}
