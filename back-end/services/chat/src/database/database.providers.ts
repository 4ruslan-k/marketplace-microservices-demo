import { Logger } from '@nestjs/common';
import { ConfigService } from '@nestjs/config';
import { User } from 'src/user/user.entity';
import { DATA_SOURCE } from 'src/user/user.providers';
import { DataSource } from 'typeorm';

export const databaseProviders = [
  {
    provide: DATA_SOURCE,
    inject: [ConfigService],
    useFactory: async (configService: ConfigService) => {
      const dataSource = new DataSource({
        type: 'postgres',
        host: configService.get('database.host'),
        port: configService.get('database.port'),
        username: configService.get('database.username'),
        password: configService.get('database.password'),
        database: configService.get('database.name'),
        logging: configService.get('database.debugMode'),
        entities: [User],
        synchronize: false,
      });

      Logger.log('Connecting to the database...');

      const source = await dataSource.initialize();

      Logger.log(
        `Successfully Connected to the database. Name: ${dataSource.options.database}`,
      );

      return source;
    },
  },
];
