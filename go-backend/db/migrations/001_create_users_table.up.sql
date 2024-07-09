-- Create users table
CREATE TABLE IF NOT EXISTS users (
    user_id SERIAL PRIMARY KEY,
    user_name VARCHAR(50) UNIQUE NOT NULL,
    first_name VARCHAR(255) NOT NULL,
    last_name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    user_status VARCHAR(1) NOT NULL CHECK (user_status IN ('I', 'A', 'T')),
    department VARCHAR(255) NOT NULL
);

-- Insert sample data
INSERT INTO users (user_name, first_name, last_name, email, user_status, department) VALUES 
('jdoe', 'John', 'Doe', 'john.doe@example.com', 'A', 'Engineering'),
('asmith', 'Alice', 'Smith', 'alice.smith@example.com', 'I', 'Marketing'),
('bwilson', 'Bob', 'Wilson', 'bob.wilson@example.com', 'A', 'Sales'),
('mjones', 'Mary', 'Jones', 'mary.jones@example.com', 'I', 'HR'),
('djames', 'David', 'James', 'david.james@example.com', 'T', 'Engineering'),
('lwhite', 'Linda', 'White', 'linda.white@example.com', 'A', 'Finance'),
('cgreen', 'Chris', 'Green', 'chris.green@example.com', 'T', 'Support'),
('rblack', 'Rachel', 'Black', 'rachel.black@example.com', 'I', 'Development'),
('tjohnson', 'Tom', 'Johnson', 'tom.johnson@example.com', 'A', 'Design'),
('pclark', 'Peter', 'Clark', 'peter.clark@example.com', 'T', 'Management');
