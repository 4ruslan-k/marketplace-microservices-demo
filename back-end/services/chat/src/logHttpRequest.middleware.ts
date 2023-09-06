import { Request, Response, NextFunction } from 'express';
import { Injectable, NestMiddleware, Logger } from '@nestjs/common';

@Injectable()
export class LoggerMiddleware implements NestMiddleware {
  private logger = new Logger('HTTP');

  use(request: Request, response: Response, next: NextFunction): void {
    const startAt = process.hrtime();
    const { ip, method, originalUrl } = request;
    const userAgent = request.get('user-agent') || '';

    response.on('finish', () => {
      const { statusCode } = response;
      const requestDuration = process.hrtime(startAt);
      const seconds = requestDuration[0];
      const milliseconds = requestDuration[1] / 1000000; // nanoseconds to milliseconds
      const elapsedTime = `${seconds}s ${milliseconds}ms`;
      this.logger.log(
        `${method} ${originalUrl} ${statusCode} ${elapsedTime} - ${userAgent} IP: ${ip}`,
      );
    });

    next();
  }
}
