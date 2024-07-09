import { ComponentFixture, TestBed } from '@angular/core/testing';
import { ReactiveFormsModule } from '@angular/forms';
import { MatInputModule } from '@angular/material/input';
import { MatButtonModule } from '@angular/material/button';
import { ApiService } from '../../services/api.service';
import { UserCreateComponent } from './user-create.component';
import { of } from 'rxjs';

describe('UserCreateComponent', () => {
  let component: UserCreateComponent;
  let fixture: ComponentFixture<UserCreateComponent>;
  let apiService: ApiService;

  beforeEach(async () => {
    const apiServiceStub = {
      createUser: jasmine.createSpy('createUser').and.returnValue(of({}))
    };

    await TestBed.configureTestingModule({
      declarations: [UserCreateComponent],
      imports: [ReactiveFormsModule, MatInputModule, MatButtonModule],
      providers: [{ provide: ApiService, useValue: apiServiceStub }]
    }).compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(UserCreateComponent);
    component = fixture.componentInstance;
    apiService = TestBed.inject(ApiService);
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  it('should call createUser on ApiService when form is valid and submitted', () => {
    component.userForm.setValue({
      userName: 'testuser',
      firstName: 'Test',
      lastName: 'User',
      email: 'testuser@example.com',
      userStatus: 'A',
      department: 'Engineering'
    });

    component.onSubmit();

    expect(apiService.createUser).toHaveBeenCalled();
  });
});
