import { Controller, Logger } from '@nestjs/common';
import { Ctx, NatsContext, Payload, EventPattern } from '@nestjs/microservices';
import { CreateUserDto, UserService } from './user.service';

@Controller()
export class UserController {
  private logger = new Logger('UserController');
  constructor(private readonly userService: UserService) {}

  @EventPattern('users.created')
  async createUser(
    @Payload() user: CreateUserDto,
    @Ctx() context: NatsContext,
  ) {
    this.logger.log({ data: user, subject: context.getSubject() });
    await this.userService.createUser(user);
  }
}
