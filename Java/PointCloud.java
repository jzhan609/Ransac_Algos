package Project;

import java.io.BufferedWriter;
import java.io.FileReader;
import java.io.FileWriter;
import java.io.IOException;
import java.util.ArrayList;
import java.util.Iterator;
import java.util.List;
import java.util.Scanner;

public class PointCloud implements Iterable<Point3D> {

    private List<Point3D> points = new ArrayList<>();

    public PointCloud(){}

    public PointCloud(String filePath) throws IOException {

        FileReader fileReader = new FileReader(filePath);
        points = new ArrayList<>();
      
        try (Scanner sc = new Scanner(fileReader)) {

            sc.nextLine();

            while (sc.hasNextLine()) {
                
                String line = sc.nextLine();
                String[] column = line.split("\t");
                double x = Double.parseDouble(column[0]);
                double y = Double.parseDouble(column[1]);
                double z = Double.parseDouble(column[2]);
                Point3D point = new Point3D(x, y, z);
                points.add(point);
            }
        }
    }

    public Point3D getPoint() {

        int r = (int)(Math.random()*points.size());
        return points.get(r);
    }

    public void addPoint(Point3D p){

        points.add(p);
    }

    public void save(String file) {

        try {

          BufferedWriter writer = new BufferedWriter(new FileWriter(file));
          
          writer.write("x\ty\tz");
          writer.newLine();

          for (Point3D p : points) {

            writer.write(p.getX() + " " + p.getY() + " " + p.getZ());
            writer.newLine();
          }
          writer.close();
        } catch (IOException e) {

            System.out.println("Oops");
        }
    }

    public Iterator<Point3D> iterator() {
        
        return points.iterator();
    }

    public int size(){ return points.size(); }
}
