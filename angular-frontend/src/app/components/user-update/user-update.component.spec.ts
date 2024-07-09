import { ComponentFixture, TestBed } from '@angular/core/testing';
import { ReactiveFormsModule } from '@angular/forms';
import { MatInputModule } from '@angular/material/input';
import { MatButtonModule } from '@angular/material/button';
import { ApiService } from '../../services/api.service';
import { UserUpdateComponent } from './user-update.component';
import { of } from 'rxjs';

describe('UserUpdateComponent', () => {
  let component: UserUpdateComponent;
  let fixture: ComponentFixture<UserUpdateComponent>;
  let apiService: ApiService;

  beforeEach(async () => {
    const apiServiceStub = {
      updateUser: jasmine.createSpy('updateUser').and.returnValue(of({}))
    };

    await TestBed.configureTestingModule({
      declarations: [UserUpdateComponent],
      imports: [ReactiveFormsModule, MatInputModule, MatButtonModule],
      providers: [{ provide: ApiService, useValue: apiServiceStub }]
    }).compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(UserUpdateComponent);
    component = fixture.componentInstance;
    apiService = TestBed.inject(ApiService);
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  it('should call updateUser on ApiService when form is valid and submitted', () => {
    component.userForm.setValue({
      userName: 'testuser',
      firstName: 'Test',
      lastName: 'User',
      email: 'testuser@example.com',
      userStatus: 'A',
      department: 'Engineering'
    });

    component.onSubmit();

    expect(apiService.updateUser).toHaveBeenCalled();
  });
});
