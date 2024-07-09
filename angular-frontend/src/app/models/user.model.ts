export interface User {
    user_id: number;
    user_name: string;
    first_name: string;
    last_name: string;
    email: string;
    user_status: 'A' | 'I' | 'T';
    department: string;
  }
  