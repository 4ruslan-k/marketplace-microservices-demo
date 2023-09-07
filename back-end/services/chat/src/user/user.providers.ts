import { DataSource } from 'typeorm';
import { User } from './user.entity';

export const USER_REPOSITORY = 'USER_REPOSITORY';
export const DATA_SOURCE = 'DATA_SOURCE';

export const userProviders = [
  {
    provide: USER_REPOSITORY,
    useFactory: (dataSource: DataSource) => dataSource.getRepository(User),
    inject: [DATA_SOURCE],
  },
];
