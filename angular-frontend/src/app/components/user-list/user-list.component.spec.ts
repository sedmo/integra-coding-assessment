import { ComponentFixture, TestBed } from '@angular/core/testing';
import { MatTableModule } from '@angular/material/table';
import { ApiService } from '../../services/api.service';
import { UserListComponent } from './user-list.component';
import { of } from 'rxjs';

describe('UserListComponent', () => {
  let component: UserListComponent;
  let fixture: ComponentFixture<UserListComponent>;
  let apiService: ApiService;

  beforeEach(async () => {
    const apiServiceStub = {
      getUsers: jasmine.createSpy('getUsers').and.returnValue(of([]))
    };

    await TestBed.configureTestingModule({
      declarations: [UserListComponent],
      imports: [MatTableModule],
      providers: [{ provide: ApiService, useValue: apiServiceStub }]
    }).compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(UserListComponent);
    component = fixture.componentInstance;
    apiService = TestBed.inject(ApiService);
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  it('should call getUsers on ApiService when initialized', () => {
    expect(apiService.getUsers).toHaveBeenCalled();
  });
});
