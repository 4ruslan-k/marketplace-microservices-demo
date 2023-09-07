import { Injectable, Inject } from '@nestjs/common';
import { Repository } from 'typeorm';
import { User } from './user.entity';
import { USER_REPOSITORY } from './user.providers';

@Injectable()
export class UserService {
  constructor(
    @Inject(USER_REPOSITORY)
    private userRepository: Repository<User>,
  ) {}

  async createUser(user: CreateUserDto) {
    await this.userRepository.insert({
      id: user.id,
      name: user.name,
      email: user.email,
      createdAt: new Date(),
      updated_at: null,
    });
  }
}

export class CreateUserDto {
  id: string;
  name: string;
  email: string;
}
