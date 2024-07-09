import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { environment } from '../../environments/environment';
import { Observable } from 'rxjs';
import { User } from '../models/user.model'; // Import the User type from the models folder

@Injectable({
  providedIn: 'root'
})
export class ApiService {
  private apiUrl = environment.apiUrl;

  constructor(private http: HttpClient) {}

  getUsers(): Observable<User[]> {
    return this.http.get<User[]>(`${this.apiUrl}/users`); // Specify the type of the response as User[]
  }

  getUserById(id: number) {
    return this.http.get(`${this.apiUrl}/users/${id}`);
  }

  createUser(user: any) {
    return this.http.post(`${this.apiUrl}/users`, user);
  }

  updateUser(id: number, user: any) {
    return this.http.put(`${this.apiUrl}/users/${id}`, user);
  }

  deleteUser(id: number) {
    return this.http.delete(`${this.apiUrl}/users/${id}`);
  }
  getApiUrl() {
    return this.apiUrl;
  }
}