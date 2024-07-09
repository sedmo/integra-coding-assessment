import { Component, OnInit } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { ApiService } from '../../services/api.service';

const VALID_USER_STATUSES = ['A', 'I', 'T'];

@Component({
  selector: 'app-user-update',
  templateUrl: './user-update.component.html',
  styleUrls: ['./user-update.component.css']
})
export class UserUpdateComponent implements OnInit {
  userId!: number;
  userForm: FormGroup;
  errorMessage: string | null = null;

  constructor(
    private route: ActivatedRoute,
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

  ngOnInit(): void {
    this.userId = this.route.snapshot.params['id'];
    this.apiService.getUserById(this.userId).subscribe(user => {
      this.userForm.patchValue(user);
    });
  }

  onSubmit(): void {
    if (this.userForm.valid) {
      this.apiService.updateUser(this.userId, this.userForm.value).subscribe(
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
