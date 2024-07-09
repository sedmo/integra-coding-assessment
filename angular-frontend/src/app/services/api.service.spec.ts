import { TestBed } from '@angular/core/testing';
import { HttpClientTestingModule, HttpTestingController } from '@angular/common/http/testing';
import { ApiService } from './api.service';
import { User } from '../models/user.model'; // Import the User type from the models folder

describe('ApiService', () => {
  let service: ApiService;
  let httpMock: HttpTestingController;

  beforeEach(() => {
    TestBed.configureTestingModule({
      imports: [HttpClientTestingModule],
      providers: [ApiService]
    });

    service = TestBed.inject(ApiService);
    httpMock = TestBed.inject(HttpTestingController);
  });

  afterEach(() => {
    httpMock.verify();
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });

  it('should fetch users', () => {
    const dummyUsers: User[] = [
      { user_id: 100, user_name: 'testuser1', first_name: 'Test', last_name: 'User', email: 'testuser1@example.com', user_status: 'A', department: 'Engineering' },
      { user_id: 101, user_name: 'testuser2', first_name: 'Test', last_name: 'User', email: 'testuser2@example.com', user_status: 'I', department: 'Marketing' }
    ];

    service.getUsers().subscribe(users => {
      expect(users.length).toBe(2);
      expect(users).toEqual(dummyUsers);
    });

    const req = httpMock.expectOne(`${service.getApiUrl()}/users`);
    expect(req.request.method).toBe('GET');
    req.flush(dummyUsers);
  });

  it('should create a user', () => {
    const newUser: User = { user_id: 102, user_name: 'testuser3', first_name: 'New', last_name: 'User', email: 'testuser3@example.com', user_status: 'T', department: 'HR' };

    service.createUser(newUser).subscribe(user => {
      expect(user).toEqual(newUser);
    });

    const req = httpMock.expectOne(`${service.getApiUrl()}/users`);
    expect(req.request.method).toBe('POST');
    req.flush(newUser);
  });

  it('should update a user', () => {
    const updatedUser: User = { user_id: 100, user_name: 'testuser1', first_name: 'Updated', last_name: 'User', email: 'updateduser@example.com', user_status: 'A', department: 'Engineering' };

    service.updateUser(updatedUser.user_id, updatedUser.user_id).subscribe(user => {
      expect(user).toEqual(updatedUser);
    });

    const req = httpMock.expectOne(`${service.getApiUrl()}/users/${updatedUser.user_id}`);
    expect(req.request.method).toBe('PUT');
    req.flush(updatedUser);
  });

  it('should delete a user', () => {
    const userId = 100;

    service.deleteUser(userId).subscribe(response => {
      expect(response).toEqual({});
    });

    const req = httpMock.expectOne(`${service.getApiUrl()}/users/${userId}`);
    expect(req.request.method).toBe('DELETE');
    req.flush({});
  });
});
