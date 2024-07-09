import { Component, OnInit } from '@angular/core';
import { Router } from '@angular/router';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { ApiService } from '../../services/api.service';

const VALID_USER_STATUSES = ['A', 'I', 'T'];

@Component({
  selector: 'app-user-create',
  templateUrl: './user-create.component.html',
  styleUrls: ['./user-create.component.css']
})
export class UserCreateComponent implements OnInit {
  userForm: FormGroup;
  errorMessage: string | null = null;

  constructor(
    private router: Router,
    private fb: FormBuilder,
    private apiService: ApiService
  ) {
    this.userForm = this.fb.group({
      user_name: ['', [Validators.required, Validators.maxLength(50)]],
      first_name: ['', [Validators.required, Validators.maxLength(50)]],
      last_name: ['', [Validators.required, Validators.maxLength(50)]],
      email: ['', [Validators.required, Validators.email, Validators.maxLength(100)]],
      user_status: ['', [Validators.required, Validators.pattern(VALID_USER_STATUSES.join('|'))]],
      department: ['', [Validators.required, Validators.maxLength(50)]]
    });
  }

  ngOnInit(): void {}

  onSubmit(): void {
    if (this.userForm.valid) {
      this.apiService.createUser(this.userForm.value).subscribe(
        () => {
          this.router.navigate(['/users']);
        },
        (error) => {
          this.errorMessage = error.error;
        }
      );
    }
  }
}
